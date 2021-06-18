Feature: Authorize Order

  Scenario: Orders are authorized on active accounts
    Given I create an account with:
      | ConsumerID | a123        |
      | Name       | TestAccount |
    When I authorize an order with:
      | ConsumerID | a123 |
      | OrderID    | b456 |
      | OrderTotal | 999  |
    Then I expect the command to succeed

  Scenario: Orders are not authorized on inactive accounts
    Given I create an account with:
      | ConsumerID | a123        |
      | Name       | TestAccount |
    And I disable the account with:
      | AccountID | a123 |
    When I authorize an order with:
      | ConsumerID | a123 |
      | OrderID    | b456 |
      | OrderTotal | 999  |
    Then I expect the command to fail
    And the returned error message is:
    """
    account is disabled
    """
