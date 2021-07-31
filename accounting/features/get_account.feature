@query
Feature: Get Accounts

  Scenario: Get an account by ID
    Given I create an account for the consumer "Able Anders"
    When I request the account for "Able Anders"
    Then I expect the request to succeed
    And the returned account is enabled

  Scenario: Get a disabled account by ID
    Given I create an account for the consumer "Able Anders"
    And I disable the account for "Able Anders"
    When I request the account for "Able Anders"
    Then I expect the request to succeed
    And the returned account is disabled

  Scenario: Getting an account that does not exist returns an error
    Given I create an account for the consumer "Able Anders"
    When I request the account for "Betty Burns"
    Then I expect the request to fail
    And the returned error message is "account not found"
