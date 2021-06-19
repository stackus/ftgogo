@command @consumer @order
Feature: Validate Orders By Consumer

  Background: Setup a consumer
    Given I setup a consumer with:
    """
    {
      "Name": "TestConsumer"
    }
    """

  Scenario: Can validate orders for consumers
    When I validate an order for the consumer with:
    """
    {
      "ConsumerID": "<ConsumerID>",
      "OrderID": "a123",
      "OrderTotal": 999
    }
    """
    Then I expect the command to succeed

  Scenario: Cannot validate orders for consumers that do not exist
    When I validate an order for the consumer with:
    """
    {
      "ConsumerID": "b456",
      "OrderID": "a123",
      "OrderTotal": 999
    }
    """
    Then I expect the command to fail
    And the returned error message is:
    """
    consumer not found
    """
