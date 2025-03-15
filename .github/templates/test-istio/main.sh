cd ${PROJECT_PATH}/99-labs/code/splitdim/deploy || exit 1
kubectl apply -f kubernetes-istio.yaml
bash ${ACTION_PATH}/wait.sh
export EXTERNAL_IP=$(kubectl get service splitdim-istio -o jsonpath='{.status.loadBalancer.ingress[0].ip}')
export EXTERNAL_PORT=80
cd ..
go test ./... --tags=httphandler,api,localconstructor,reset,transfer,accounts,clear -v -count 1