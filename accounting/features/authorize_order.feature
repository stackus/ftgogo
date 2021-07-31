@command
Feature: Authorize Order

  Scenario: Orders are authorized on active accounts
    Given I create an account for the consumer "Able Anders"
    When I authorize an order totaling $9.99 for "Able Anders"
    Then I expect it be authorized

  Scenario: Orders are not authorized on inactive accounts
    Given I create an account for the consumer "Able Anders"
    And I disable the account for "Able Anders"
    When I authorize an order totaling $9.99 for "Able Anders"
    Then I don't expect it to be authorized
    And the returned error message is "account is disabled"

  Scenario: Orders are not authorized on unregistered accounts
    When I authorize an order totaling $9.99 for "Able Anders"
    Then I don't expect it to be authorized
    And the returned error message is "account not found"
