SECONDS=0
while true; do
  kubectl get svc #for debug
  IP=$(kubectl get service splitdim-istio -o jsonpath='{.status.loadBalancer.ingress[0].ip}')
  if [[ ! -n "$IP" ]]; then
    if (( SECONDS == 15));then
      exit 1
    fi
    SECONDS=$(( SECONDS + 1 ))
    sleep 1
  else
    exit 0
  fi
done