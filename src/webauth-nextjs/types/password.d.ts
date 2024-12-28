export declare type IProviders = "google" | "facebook";

export interface IAPIError {
  name: string;
  message: string;
}

export interface IPasswordRules {
  label: string;
  pattern: string;
  checked?: boolean;
}

export interface PasswordCheckListProps {
  className?: string;
  value: string;
  onChange: (valid: boolean) => void;
  rules?: IPasswordRules[];
}

export interface IDisplayErrorProps {
  name?: string;
  message: string;
}
