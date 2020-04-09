#!/bin/bash
sleep 5
curl -X POST http://localhost:$DAPR_HTTP_PORT/v1.0/state/statestore \
  -H "Content-Type: application/json" \
  -d '[
        {
          "key": "1",
          "value": {
            
          }
        },
        {
          "key": "2",
          "value": {
            "name": "Lemon Curd"
          }
        }
      ]'