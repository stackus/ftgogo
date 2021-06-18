Feature: Enable Accounts

  Scenario: Disabled accounts can be re-enabled
    Given I create an account with:
      | ConsumerID | a123        |
      | Name       | TestAccount |
    And I disable the account with:
      | AccountID | a123 |
    When I enable the account with:
      | AccountID | a123 |
    Then I expect the command to succeed

  Scenario: Enabling already enabled accounts return an error
    Given I create an account with:
      | ConsumerID | a123        |
      | Name       | TestAccount |
    When I enable the account with:
      | AccountID | a123 |
    Then I expect the command to fail
    And the returned error message is:
    """
    account is enabled
    """

  Scenario: Enabling accounts that do not exist returns an error
    Given I create an account with:
      | ConsumerID | a123        |
      | Name       | TestAccount |
    And I disable the account with:
      | AccountID | a123 |
    When I enable an account with:
      | AccountID | b456 |
    Then I expect the command to fail
    And the returned error message is:
    """
    account not found
    """
