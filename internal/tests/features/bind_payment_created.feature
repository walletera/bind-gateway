Feature: process Bind webhook event transfer.cvu.received
  Bind sends a webhook event of type transfer.cvu.received

  Background: the bind-gateway is up and running
    Given a running bind-gateway

  Scenario: a bind inbound transfer is processed successfully
    Given a Bind transfer.cvu.received event:
    """json
    {
      "id": "9f7e2845-2ca3-4787-826a-849f94ba469f",
      "object": "",
      "created": "2023-06-07T18:22:16.0498297",
      "data": {
        "id": "1-30717449076-000000005536667-1",
        "type": "TRANSFER",
        "from": {
          "bank_id": "322",
          "account_id": "20-1-735135-1-5"
        },
        "counterparty": {
          "id": "20322204121",
          "name": "Perez Juan",
          "id_type": "CUIT_CUIL",
          "bank_routing": {
            "scheme": null,
            "address": ""
          },
          "account_routing": {
            "scheme": "CVU",
            "address": "0000254900000000201403"
          }
        },
        "details": {
          "origin_id": 1234567,
          "origin_debit": {
            "cvu": "0000547800000000201970",
            "cuit": 20110006668
          },
          "origin_credit": {
            "cvu": "0000347800000000201874",
            "cuit": 23112223339
          }
        },
        "transaction_ids": [
          "1-30717449076-000000005536667-1",
          "J3D5W612E7754YQNGXYVRL"
        ],
        "status": "COMPLETED",
        "start_date": "2023-06-07T18:22:17",
        "end_date": "2023-06-07T18:22:17",
        "challenge": null,
        "charge": {
          "summary": "VAR primary",
          "value": {
            "currency": "ARS",
            "amount": 80
          }
        }
      },
      "type": "transfer.cvu.received",
      "redeliveries": 0
    }
    """
    And  an accounts endpoint to get accounts:
    # the json below is a mockserver expectation
    """json
    {
      "id": "getAccountSucceed",
      "httpRequest" : {
        "method": "GET",
        "path": "/accounts",
        "queryStringParameters": {
            "accountType": ["cvu"],
            "cvuAccountDetails[cvu]": ["0000347800000000201874"],
            "cvuAccountDetails[cuit]": ["23112223339"]
        }
      },
      "httpResponse" : {
        "statusCode" : 200,
        "headers" : {
          "content-type" : [ "application/json" ]
        },
        "body": [{
            "id": "01937863-163a-790f-8e59-707e152dd9c7",
            "currency": "USD",
            "customerId": "9fd3bc09-99da-4486-950a-11082f5fd966",
            "customerName": "Lemon Cash",
            "customerAccountId": "UID76ASDF87",
            "accountType": "cvu",
            "accountDetails": {
                "cvu": "0000347800000000201874",
                "cuit": "23112223339"
            }
        }]
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
    And  a payments endpoint to create payments:
    # the json below is a mockserver expectation
    """json
    {
      "id": "postDepositSucceed",
      "httpRequest" : {
        "method": "POST",
        "path": "/payments",
        "body": {
            "type": "JSON",
            "json": {
              "id": "${json-unit.any-string}",
              "amount": 80,
              "currency": "ARS",
              "customerId": "9fd3bc09-99da-4486-950a-11082f5fd966",
              "externalId": "1234567",
              "beneficiary": {
                "currency": "ARS",
                "accountType": "cvu",
                "accountDetails": {
                    "cvu": "0000347800000000201874",
                    "cuit": "23112223339"
                }
              },
              "debtor": {
                "currency": "ARS",
                "accountType": "cvu",
                "accountDetails": {
                    "cvu": "0000547800000000201970",
                    "cuit": "20110006668"
                }
              }
              },
            "matchType": "ONLY_MATCHING_FIELDS"
        }
      },
      "httpResponse" : {
        "statusCode" : 201,
        "headers" : {
          "content-type" : [ "application/json" ]
        },
        "body": {
            "id": "01937863-163a-790f-8e59-707e152dd9c7",
            "amount": 100,
            "currency": "USD",
            "direction": "outbound",
            "customerId": "9fd3bc09-99da-4486-950a-11082f5fd966",
            "externalId": "1234567",
            "status": "confirmed",
            "createdAt": "2024-11-29T14:46:19Z"
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
    When the webhook event is received
    Then the bind-gateway creates the corresponding payment on the Payments API
    And the bind-gateway produces the following log:
    """
    bind event InboundTransferCreated processed successfully
    """
    And the bind-gateway produces the following log:
    """
    gateway event InboundPaymentReceived processed successfully
    """
