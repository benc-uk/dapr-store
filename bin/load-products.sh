POD_NAME=$(kubectl get pods|grep products|awk '{print $1}')

kubectl cp ./etc/load-data.sh -c service $POD_NAME:.
kubectl cp ./etc/product-data.json -c service $POD_NAME:.
kubectl exec -it $POD_NAME -c service -- sh ./load-data.sh product-data.json
