Feature: process PaymentCreated event
  Walletera users with funds on their accounts want to be able to Withdraw those funds (totally or partially)
  to DinoPay accounts.
  - When a Withdrawal is created, A PaymentCreated event is published by the Payments Service.
  - The DinoPay Gateway must listen for PaymentCreated events.
  - Whenever a PaymentCreated event arrives the bind-gateway must create a corresponding Payment on DinoPay API.

  Background: the bind-gateway is up and running
    Given a running bind-gateway

  Scenario: payment created event is processed successfully
    Given a PaymentCreated event:
    """json
    {
      "id": "0fg1833e-3438-4908-b90a-5721670cb067",
      "type": "PaymentCreated",
      "data": {
        "id": "0ae1733e-7538-4908-b90a-5721670cb093",
        "amount": 100,
        "currency": "USD",
        "direction": "outbound",
        "customerId": "2432318c-4ff3-4ac0-b734-9b61779e2e46",
        "status": "pending",
        "debtor": {
          "bankName": "LetsBit",
          "bankId": "letsbit",
          "accountHolder": "Ron Doe",
          "routingKey": "0003252627188236545234"
        },
        "beneficiary": {
          "bankName": "LetsBit",
          "bankId": "letsbit",
          "accountHolder": "Richard Roe",
          "routingKey": "0004252627182736545234"
        },
        "createdAt": "2024-10-04T00:00:00Z"
      }
    }
    """
    And  a dinopay endpoint to create payments:
    # the json below is a mockserver expectation
    """json
    {
      "id": "createPaymentSucceed",
      "httpRequest" : {
        "method": "POST",
        "path" : "/walletentidad-operaciones/v1/api/v1.201/transferir",
        "body": {
            "type": "JSON",
            "json": {
              "cvuOrigen": "0003252627188236545234",
              "cbu_cvu_destino": "0004252627182736545234",
              "importe": 100,
              "idExterno": "0ae1733e-7538-4908-b90a-5721670cb093"
            }
        }
      },
      "httpResponse" : {
        "statusCode" : 200,
        "headers" : {
          "content-type" : [ "application/json" ]
        },
        "body" : {
          "operacionId": 584866,
          "operacionIdExterno": "0fg1833e-3438-4908-b90a-5721670cb067",
          "estadoExterno": "COMPLETED",
          "estadoId": 2,
          "origenCuentaId": 274931,
          "coelsaId": "KLOEJWV9JR48G18NQMD0GZ",
          "fechaInicio": "2024-09-03T23:29:29+00:00",
          "fechaFin": "2024-09-03T23:29:29+00:00",
          "fechaNegocio": "2024-09-04T03:00:00+00:00",
          "importe": 100,
          "cvuOrigen": "0003252627188236545234",
          "referencia": "futbol 5",
          "concepto": "VAR",
          "cvuCbuContraparte": "0004252627182736545234",
          "nombreContraparte": "GRANJAS CARNAVE SA",
          "cuitCuilContraparte": "30707101020",
          "comprobanteId": 8470933,
          "esTransferenciaInterna": false,
          "estaFinalizada": true,
          "estaRechazada": false,
          "estaAAuditar": false,
          "estaPendiente": false
        }
      },
      "priority" : 0,
      "timeToLive" : {
        "unlimited" : true
      },
      "times" : {
        "unlimited" : true
      }
    }
    """
    And  a payments endpoint to update payments:
    # the json below is a mockserver expectation
    """json
    {
      "id": "updatePaymentSucceed",
      "httpRequest" : {
        "method": "PATCH",
        "path": "/payments/0ae1733e-7538-4908-b90a-5721670cb093",
        "body": {
          "type": "JSON",
          "json": {
            "externalId": "584866",
            "status": "confirmed"
          }
        }
      },
      "httpResponse" : {
        "statusCode" : 200,
        "headers" : {
          "content-type" : [ "application/json" ]
        }
      },
      "priority" : 0,
      "timeToLive" : {
        "unlimited" : true
      },
      "times" : {
        "unlimited" : true
      }
    }
    """
    When the event is published
    Then the bind-gateway creates the corresponding payment on the Bind API
    And the bind-gateway updates the payment on payments service
    And the bind-gateway produces the following log:
    """
    PaymentCreated event processed successfully
    """
    And the bind-gateway produces the following log:
    """
    OutboundPaymentCreated event processed successfully
    """

#  Scenario: payment created event processing failed when trying to create payment on Dinopay
#    Given a PaymentCreated event:
#    """json
#    {
#      "id": "0fg1833e-3438-4908-b90a-5721670cb067",
#      "type": "PaymentCreated",
#      "data": {
#        "id": "0ae1733e-7538-4908-b90a-5721670cb093",
#        "amount": 100,
#        "currency": "USD",
#        "direction": "outbound",
#        "customerId": "2432318c-4ff3-4ac0-b734-9b61779e2e46",
#        "status": "pending",
#        "beneficiary": {
#          "bankName": "dinopay",
#          "bankId": "dinopay",
#          "accountHolder": "Richard Roe",
#          "routingKey": "123456789123456",
#          "accountNumber": "1200079635"
#        },
#        "createdAt": "2024-10-04T00:00:00Z"
#      }
#    }
#    """
#    And  a dinopay endpoint to create payments:
#    # the json below is a mockserver expectation
#    """json
#    {
#      "id": "createPaymentFail",
#      "httpRequest" : {
#        "method": "POST",
#        "path" : "/payments",
#        "body": {
#            "type": "JSON",
#            "json": {
#              "customerTransactionId": "0ae1733e-7538-4908-b90a-5721670cb093",
#              "amount": 100,
#              "currency": "USD",
#              "destinationAccount": {
#                "accountHolder": "Richard Roe",
#                "accountNumber": "1200079635"
#              }
#            },
#            "matchType": "ONLY_MATCHING_FIELDS"
#        }
#      },
#      "httpResponse" : {
#        "statusCode" : 500,
#        "headers" : {
#          "content-type" : [ "text/html" ]
#        },
#        "body" : "something bad happened"
#      },
#      "priority" : 0,
#      "timeToLive" : {
#        "unlimited" : true
#      },
#      "times" : {
#        "unlimited" : true
#      }
#    }
#    """
#    When the event is published
#    Then the bind-gateway fails creating the corresponding payment on the DinoPay API
#    And  the bind-gateway produces the following log:
#    """
#    failed creating payment on dinopay
#    """
