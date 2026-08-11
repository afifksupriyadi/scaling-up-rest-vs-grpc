import grpc from 'k6/net/grpc';

const client = new grpc.Client();
client.load(['/home/afiffks/scaling-up-rest-vs-grpc/proto'], 'order.proto');

const METHOD = __ENV.GRPC_METHOD;
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
  if (__ITER == 0) {
    client.connect('10.184.0.2:50051', { plaintext: true });
  }
  client.invoke(`orderexperiment.OrderExperimentService/${METHOD}`, {});
};