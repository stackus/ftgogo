@query @consumer
Feature: Get Consumer

  Background: Setup a consumer
    Given I setup a consumer with:
    """
    {
      "Name": "TestConsumer"
    }
    """

  Scenario: Can get consumers
    When I request a consumer with:
    """
    {
      "ConsumerID": "<ConsumerID>"
    }
    """
    Then I expect the request to succeed

  Scenario: Asking for a consumer that does not exist returns an error
    When I request a consumer with:
    """
    {
      "ConsumerID": "b456"
    }
    """
    Then I expect the request to fail
    And the returned error message is:
    """
    consumer not found
    """
