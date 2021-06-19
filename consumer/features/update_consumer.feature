@command @consumer
Feature: Update Consumers

  Background: Setup a consumer
    Given I setup a consumer with:
    """
    {
      "Name": "TestConsumer"
    }
    """

  Scenario: Consumers can be updated
    When I update the consumer with:
    """
    {
      "ConsumerID": "<ConsumerID>",
      "Name": "UpdatedConsumer"
    }
    """
    Then I expect the command to succeed

  Scenario: Updating consumers that do not exist returns an error
    When I update a consumer with:
    """
    {
      "ConsumerID": "b456",
      "Name": "OtherConsumer"
    }
    """
    Then I expect the command to fail
    And the returned error message is:
    """
    consumer not found
    """
