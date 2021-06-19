@command
Feature: Disable Accounts

  Scenario: Enabled accounts can be disabled
    Given I create an account with:
    """
    {
      "ConsumerID": "a123",
      "Name": "TestAccount"
    }
    """
    When I disable the account with:
    """
    {
      "AccountID": "a123"
    }
    """
    Then I expect the command to succeed

  Scenario: Disabling already disabled accounts return an error
    Given I create an account with:
    """
    {
      "ConsumerID": "a123",
      "Name": "TestAccount"
    }
    """
    And I disable the account with:
    """
    {
      "AccountID": "a123"
    }
    """
    When I disable the account with:
    """
    {
      "AccountID": "a123"
    }
    """
    Then I expect the command to fail
    And the returned error message is:
    """
    account is disabled
    """

  Scenario: Disabling accounts that do not exist returns an error
    Given I create an account with:
    """
    {
      "ConsumerID": "a123",
      "Name": "TestAccount"
    }
    """
    When I disable an account with:
    """
    {
      "AccountID": "b456"
    }
    """
    Then I expect the command to fail
    And the returned error message is:
    """
    account not found
    """
