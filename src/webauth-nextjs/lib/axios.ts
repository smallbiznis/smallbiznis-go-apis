import axios from "axios";

export const api = axios.create({
  baseURL: process.env.APP_URL || 'https://accounts.smallbiznis.test',
  insecureHTTPParser: true,
  validateStatus: (status) => {
    return status >= 200 && status < 500;
  },
});
