Feature: Get Accounts

  Scenario: Get an account by ID
    Given I create an account with:
      | ConsumerID | a123        |
      | Name       | TestAccount |
    When I request an account with:
      | AccountID | a123 |
    Then I expect the command to succeed
    And the returned account matches:
    """
    {
      "ID": "a123",
      "Name": "TestAccount",
      "Enabled": true
    }
    """

  Scenario: Get a disabled account by ID
    Given I create an account with:
      | ConsumerID | a123        |
      | Name       | TestAccount |
    And I disable the account with:
      | AccountID | a123 |
    When I request an account with:
      | AccountID | a123 |
    Then I expect the command to succeed
    And the returned account matches:
    """
    {
      "ID": "a123",
      "Name": "TestAccount",
      "Enabled": false
    }
    """

  Scenario: Getting an account that does not exist returns an error
    Given I create an account with:
      | ConsumerID | a123        |
      | Name       | TestAccount |
    When I request an account with:
      | AccountID | b456 |
    Then I expect the command to fail
    And the returned error message is:
    """
    account not found
    """
