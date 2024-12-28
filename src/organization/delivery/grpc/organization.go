package grpc

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/gosimple/slug"
	"github.com/smallbiznis/go-genproto/smallbiznis/balance/v1"
	"github.com/smallbiznis/go-genproto/smallbiznis/organization/v1"
	"github.com/smallbiznis/go-genproto/smallbiznis/subscription/v1"
	"github.com/smallbiznis/go-lib/pkg/env"
	"github.com/smallbiznis/go-lib/pkg/pagination"
	"github.com/smallbiznis/organization/domain"
	"github.com/stripe/stripe-go/v80"
	"github.com/stripe/stripe-go/v80/billingportal/session"
	"github.com/stripe/stripe-go/v80/customer"
	"github.com/stripe/stripe-go/v80/price"
	sub "github.com/stripe/stripe-go/v80/subscription"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"gorm.io/gorm"
)

type OrganizationServiceSever struct {
	organization.UnimplementedServiceServer
	db          *gorm.DB
	balanceConn balance.BalanceServiceClient
	// memberConn  member.MemberServiceClient
	// posConn                pos.ServiceClient
	countryRepo            domain.ICountryRepository
	organizationRepo       domain.IOrganizationRepository
	locationRepo           domain.ILocationRepository
	taxRepository          domain.ITaxRulesRepository
	shippingRateRepository domain.IShippingRateRepository
}

func NewOrganizationServiceServer(
	db *gorm.DB,
	balanceConn balance.BalanceServiceClient,
	// memberConn member.MemberServiceClient,
	// posConn pos.ServiceClient,
	countryRepo domain.ICountryRepository,
	organizationRepo domain.IOrganizationRepository,
	locationRepo domain.ILocationRepository,
	taxRepository domain.ITaxRulesRepository,
	shippingRateRepository domain.IShippingRateRepository,
) *OrganizationServiceSever {
	return &OrganizationServiceSever{
		db:          db,
		balanceConn: balanceConn,
		// memberConn:  memberConn,
		// posConn:                posConn,
		countryRepo:            countryRepo,
		organizationRepo:       organizationRepo,
		locationRepo:           locationRepo,
		taxRepository:          taxRepository,
		shippingRateRepository: shippingRateRepository,
	}
}

func (srv *OrganizationServiceSever) ListOrg(ctx context.Context, req *organization.ListOrganizationRequest) (*organization.ListOrganizationResponse, error) {
	span := trace.SpanFromContext(ctx)
	defer span.End()

	span.SetName("ListOrg")

	filter := domain.Organization{}
	if req.OrganizationId != "" {
		filter.OrganizationID = req.OrganizationId
	}

	orgs, count, err := srv.organizationRepo.Find(ctx, pagination.Pagination{
		Page:    int(req.Page),
		Size:    int(req.Size),
		SortBy:  req.SortBy,
		OrderBy: req.OrderBy.String(),
	}, filter)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &organization.ListOrganizationResponse{
		TotalData: int32(count),
		Data:      orgs.ToProto(),
	}, nil
}

func (srv *OrganizationServiceSever) GetOrg(ctx context.Context, req *organization.GetOrganizationRequest) (resp *organization.Organization, err error) {
	span := trace.SpanFromContext(ctx)
	defer span.End()

	span.SetName("GetOrg")

	f := domain.Organization{}
	if _, err := uuid.Parse(req.OrganizationId); err != nil {
		f.OrganizationID = req.OrganizationId
	} else {
		f.ID = req.OrganizationId
	}

	exist, err := srv.organizationRepo.FindOne(ctx, f)
	if err != nil {
		return
	}

	if exist == nil {
		return nil, status.Error(codes.InvalidArgument, "organization not found")
	}

	return exist.ToProto(), nil
}

func (srv *OrganizationServiceSever) CreateOrg(ctx context.Context, req *organization.CreateOrganizationRequest) (resp *organization.Organization, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	span := trace.SpanFromContext(ctx)
	defer span.End()

	span.SetName("CreateOrg")

	exist, err := srv.organizationRepo.FindOne(ctx, domain.Organization{OrganizationID: slug.Make(req.Title)})
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	if exist != nil {
		return nil, status.Error(codes.InvalidArgument, "organization already exist")
	}

	country, err := srv.countryRepo.FindOne(ctx, domain.Countries{
		CountryCode: req.CountryId,
	})
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	newOrg := domain.Organization{
		ID:             uuid.NewString(),
		OrganizationID: slug.Make(req.Title),
		CountryID:      country.CountryCode,
		Title:          req.Title,
		Status:         domain.OrganizationStatus(organization.Organization_ACTIVE.String()),
	}

	if err := srv.db.WithContext(ctx).Transaction(func(tx *gorm.DB) (err error) {

		if _, err := srv.CreateTaxRule(ctx, &organization.TaxRule{
			OrganizationId: newOrg.ID,
			CountryId:      newOrg.CountryID,
			Type:           organization.TaxType_VAT,
			Rate:           0.11,
		}); err != nil {
			zap.L().Error("failed create tax rule", zap.Error(err))
			return err
		}

		if _, err := srv.CreateLocation(ctx, &organization.Location{
			OrganizationId: newOrg.ID,
			Name:           "Default Location",
			Country: &organization.Country{
				CountryCode: newOrg.CountryID,
			},
			IsDefault: true,
		}); err != nil {
			zap.L().Error("failed create default location", zap.Error(err))
			return err
		}

		// if _, err := srv.memberConn.AddMember(ctx, &member.AddMemberRequest{
		// 	UserId:        req.FirstUser.UserId,
		// 	OrganizatinId: exist.ID,
		// 	Roles: []member.Role{
		// 		member.Role_ROLE_ADMIN,
		// 	},
		// }); err != nil {
		// 	zap.L().Error("failed add first member", zap.Error(err))
		// 	return err
		// }

		// if _, err := srv.posConn.AddStaff(ctx, &pos.AddStaffRequest{
		// 	UserId:        req.FirstUser.UserId,
		// 	OrganizatinId: exist.ID,
		// 	Roles: []pos.Role{
		// 		pos.Role_ROLE_ADMIN,
		// 	},
		// }); err != nil {
		// 	zap.L().Error("failed add first staff pos", zap.Error(err))
		// 	return err
		// }

		// if _, err := srv.balanceConn.CreateBalance(ctx, &balance.CreateBalanceRequest{
		// 	OrganizationId: newOrg.ID,
		// }); err != nil {
		// 	zap.L().Error("failed create balance", zap.Error(err))
		// 	return err
		// }

		if err := tx.Create(&newOrg).Error; err != nil {
			return err
		}

		return
	}); err != nil {
		fmt.Printf("failed create organization: %v\n", err)
		return nil, status.Error(codes.Internal, err.Error())
	}

	go func(ctx context.Context, org domain.Organization) {
		srv.createSubscription(&newOrg)
	}(ctx, newOrg)

	return srv.GetOrg(ctx, &organization.GetOrganizationRequest{OrganizationId: newOrg.ID})
}

func (srv *OrganizationServiceSever) createSubscription(org *domain.Organization) (err error) {
	ctx := context.Background()

	country, err := srv.countryRepo.FindOne(ctx, domain.Countries{CountryCode: org.CountryID})
	if err != nil {
		return err
	}

	if country == nil {
		return errors.New("invalid country")
	}

	stripeCustomer, err := customer.New(&stripe.CustomerParams{
		Name: &org.Title,
	})
	if err != nil {
		return err
	}

	org.StripeCustomerID = &stripeCustomer.ID

	var trialPeriodDay int64 = 14
	days, err := strconv.Atoi(env.Lookup("STRIPE_TRIAL_PERIOD_DAYS", "14"))
	if err == nil {
		trialPeriodDay = int64(days)
	}

	priceParams := &stripe.PriceListParams{
		Active: stripe.Bool(true),
		Recurring: &stripe.PriceListRecurringParams{
			Interval: stripe.String("month"),
		},
		LookupKeys: stripe.StringSlice([]string{
			fmt.Sprintf("basic_%s", strings.ToLower(country.CurrencyCode)),
		}),
	}

	priceList := price.List(priceParams)

	for priceList.Next() {
		price := priceList.Price()

		// Create New Subscription
		stripeSubscription, err := sub.New(&stripe.SubscriptionParams{
			Customer: &stripeCustomer.ID,
			Items: []*stripe.SubscriptionItemsParams{
				{
					Price:    stripe.String(price.ID),
					Quantity: stripe.Int64(1),
				},
			},
			TrialPeriodDays: stripe.Int64(trialPeriodDay),
			TrialSettings: &stripe.SubscriptionTrialSettingsParams{
				EndBehavior: &stripe.SubscriptionTrialSettingsEndBehaviorParams{
					MissingPaymentMethod: stripe.String(string(stripe.SubscriptionTrialSettingsEndBehaviorMissingPaymentMethodPause)),
				},
			},
			PaymentSettings: &stripe.SubscriptionPaymentSettingsParams{
				PaymentMethodTypes: []*string{
					stripe.String(string(stripe.PaymentMethodTypeCard)),
				},
				SaveDefaultPaymentMethod: stripe.String(string(stripe.SubscriptionPaymentSettingsSaveDefaultPaymentMethodOnSubscription)),
			},
		})
		if err != nil {
			return err
		}

		org.StripeSubscriptionID = &stripeSubscription.ID
	}

	if err := priceList.Err(); err != nil {
		return err
	}

	if _, err := srv.organizationRepo.Update(ctx, *org); err != nil {
		return err
	}

	return
}

func (srv *OrganizationServiceSever) UpdateOrg(ctx context.Context, req *organization.Organization) (resp *organization.Organization, err error) {
	span := trace.SpanFromContext(ctx)
	defer span.End()

	span.SetName("UpdateOrg")

	return resp, nil
}

func (srv *OrganizationServiceSever) DeleteOrg(ctx context.Context, req *organization.Organization) (resp *emptypb.Empty, err error) {
	span := trace.SpanFromContext(ctx)
	defer span.End()

	span.SetName("UpdateOrg")

	if err := srv.organizationRepo.Delete(ctx, domain.Organization{
		ID: req.Id,
	}); err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return resp, nil
}

func (srv *OrganizationServiceSever) GetSubscription(ctx context.Context, req *subscription.Subscription) (resp *subscription.Subscription, err error) {
	existSub, err := sub.Get(req.Id, &stripe.SubscriptionParams{})
	if err != nil {
		return
	}

	resp = &subscription.Subscription{
		Id:                 existSub.ID,
		BillingCycleAnchor: existSub.BillingCycleAnchor,
		CancelAtPeriodEnd:  existSub.CancelAtPeriodEnd,
		CanceledAt:         existSub.CancelAt,
		Customer:           existSub.Customer.ID,
		Created:            existSub.Created,
		Currency:           string(existSub.Currency),
		CurrentPeriodStart: existSub.CurrentPeriodStart,
		CurrentPeriodEnd:   existSub.CurrentPeriodEnd,
		StartDate:          existSub.StartDate,
		Livemode:           existSub.Livemode,
		Status:             string(existSub.Status),
		LatestInvoice:      existSub.LatestInvoice.ID,
		Items: &subscription.SubscriptionItems{
			HasMore: existSub.Items.HasMore,
			Url:     existSub.Items.URL,
		},
		TrialStart: existSub.TrialStart,
		TrialEnd:   existSub.TrialEnd,
	}

	for _, v := range existSub.Items.Data {
		items := &subscription.SubscriptionItem{
			Id:           v.ID,
			Subscription: v.Subscription,
			Created:      v.Created,
			Quantity:     int32(v.Quantity),
		}

		items.Price = &subscription.Price{
			Id:            v.Price.ID,
			Active:        v.Price.Active,
			BillingScheme: string(v.Price.BillingScheme),
			Created:       v.Price.Created,
			Currency:      string(v.Price.Currency),
			Product:       v.Price.Product.ID,
			Recurring: &subscription.Price_Recurring{
				Interval:      string(v.Price.Recurring.Interval),
				IntervalCount: int32(v.Price.Recurring.IntervalCount),
				UsageType:     string(v.Price.Recurring.UsageType),
			},
			TaxBehavior:       string(v.Price.TaxBehavior),
			Type:              string(v.Price.Type),
			UnitAmount:        v.Price.UnitAmount,
			UnitAmountDecimal: fmt.Sprintf("%v", v.Price.UnitAmountDecimal),
			Livemode:          v.Price.Livemode,
			Metadata:          v.Price.Metadata,
		}

		resp.Items.Data = append(resp.Items.Data, items)
	}

	return resp, nil
}

func (srv *OrganizationServiceSever) CreateBillingPortal(ctx context.Context, req *organization.CreateBillingRequest) (resp *subscription.BillingPortalSession, err error) {

	params := &stripe.BillingPortalSessionParams{
		Customer:  stripe.String(req.CustomerId),
		ReturnURL: stripe.String(req.ReturnUrl),
		Locale:    stripe.String("en"),
	}

	if req.Locale != "" {
		params.Locale = stripe.String(req.Locale)
	}

	if req.Flow != nil {
		flow := req.Flow
		params.FlowData = &stripe.BillingPortalSessionFlowDataParams{
			Type: stripe.String(req.Flow.Type.String()),
		}

		if flow.Type == subscription.FlowType_subscription_cancel {
			subscriptionCancel := flow.SubscriptionCancel
			params.FlowData.SubscriptionCancel = &stripe.BillingPortalSessionFlowDataSubscriptionCancelParams{
				Subscription: &subscriptionCancel.SubscriptionId,
			}
		} else if flow.Type == subscription.FlowType_payment_method_update {
			subscriptionUpdate := flow.SubscriptionUpdate
			params.FlowData.AfterCompletion = &stripe.BillingPortalSessionFlowDataAfterCompletionParams{
				Type: stripe.String("portal_homepage"),
			}
			params.FlowData.SubscriptionUpdate = &stripe.BillingPortalSessionFlowDataSubscriptionUpdateParams{
				Subscription: stripe.String(subscriptionUpdate.SubscriptionId),
			}
		} else {
			return nil, status.Error(codes.InvalidArgument, fmt.Sprintf("unsupported flow type %s", flow.Type.String()))
		}
	}

	sess, err := session.New(params)
	if err != nil {
		return
	}

	var flowData *subscription.Flow
	if sess.Flow != nil {
		flow := sess.Flow
		flowData := &subscription.Flow{
			Type: subscription.FlowType(subscription.FlowType_value[string(flow.Type)]),
		}

		if flow.Type == stripe.BillingPortalSessionFlowTypeSubscriptionCancel {
			flowData.SubscriptionCancel = &subscription.FlowSubscriptionCancel{
				SubscriptionId: flow.SubscriptionCancel.Subscription,
			}
		} else if flow.Type == stripe.BillingPortalSessionFlowTypeSubscriptionUpdate {
			flowData.SubscriptionUpdate = &subscription.FlowSubscriptionUpdate{
				SubscriptionId: flow.SubscriptionUpdate.Subscription,
			}
		} else if flow.Type == stripe.BillingPortalSessionFlowTypeSubscriptionUpdateConfirm {
			flowData.SubscriptionUpdateConfirm = &subscription.FlowSubscriptionUpdateConfirm{
				Items: []*subscription.FlowSubscriptionUpdateConfirmItem{
					{
						Id:       flow.SubscriptionUpdateConfirm.Items[0].ID,
						PriceId:  flow.SubscriptionUpdateConfirm.Items[0].Price,
						Quantity: int32(flow.SubscriptionUpdateConfirm.Items[0].Quantity),
					},
				},
			}
		}
	}

	return &subscription.BillingPortalSession{
		Id:            sess.ID,
		CustomerId:    sess.Customer,
		Configuration: sess.Configuration.ID,
		Created:       int32(sess.Created),
		Flow:          flowData,
		Locale:        sess.Locale,
		Livemode:      sess.Livemode,
		ReturlUrl:     sess.ReturnURL,
		Url:           sess.URL,
	}, nil
}
