import { TracingInstrumentation } from '@grafana/faro-web-tracing';
import {
  ErrorsInstrumentation,
  initializeFaro,
  SessionInstrumentation,
  ViewInstrumentation,
  WebVitalsInstrumentation,
} from '@grafana/faro-web-sdk';

export const initSDKFaro = () => {
  return initializeFaro({
    // Mandatory, the URL of the Grafana Cloud collector with embedded application key.
    // Copy from the configuration page of your application in Grafana.
    url: process.env.FARO_AGENT_ADDR || 'http://localhost:12348/collect',
    apiKey: process.env.FARO_API_KEY || 'my_super_app_key', 
    // Mandatory, the identification label(s) of your application
    app: {
      name: process.env.SERVICE_NAME || "example",
      version: process.env.SERVICE_VERSION || "1.0.0", // Optional, but recommended,
      namespace: process.env.SERVICE_NAMESPACE || "myapp",
      environment: process.env.NODE_ENV || "development"
    },
  
    instrumentations: [
      new ErrorsInstrumentation(),
      new WebVitalsInstrumentation(),
      new TracingInstrumentation(),
      new ViewInstrumentation(),
      new SessionInstrumentation(),
    ],
  })
}