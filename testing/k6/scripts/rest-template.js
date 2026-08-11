import http from 'k6/http';

const URL = __ENV.TARGET_URL;
const VUS = parseInt(__ENV.VUS);
const ITERATIONS = parseInt(__ENV.ITERATIONS);

export const options = {
  scenarios: {
    default: {
      executor: 'per-vu-iterations',
      vus: VUS,
      iterations: ITERATIONS,
      maxDuration: '10m',
    },
  },
};

export default () => {
  http.get(URL);
};