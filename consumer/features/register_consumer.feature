@command @consumer
Feature: Register Consumer

  Scenario: Consumers can be registered
    When I register a consumer with:
    """
    {
      "Name": "TestConsumer"
    }
    """
    Then I expect the command to succeed

  Scenario: Consumers must be registered with a name
    When I register a consumer with:
    """
    {
      "Name": ""
    }
    """
    Then I expect the command to fail
    And the returned error message is:
    """
    cannot register a consumer without a name
    """

  Scenario: Duplicate consumer names do not cause conflicts
    Given I register a consumer with:
    """
    {
      "Name": "TestConsumer"
    }
    """
    When I register another consumer with:
    """
    {
      "Name": "TestConsumer"
    }
    """
    Then I expect the command to succeed
