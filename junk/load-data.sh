if [ -z $1 ]; then
  echo "Please provide JSON filename"
  exit 1
fi

sleep 2
echo "### Loading 'data.json' into state"
echo "### Dapr port is $DAPR_HTTP_PORT"
curl -X POST http://localhost:$DAPR_HTTP_PORT/v1.0/state/statestore -H "Content-Type: application/json" --data @$1
echo "### Done!"