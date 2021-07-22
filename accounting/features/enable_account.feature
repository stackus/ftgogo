@command
Feature: Enable Accounts

  Scenario: Disabled accounts can be re-enabled
    Given I create an account for the consumer "Able Anders"
    And I disable the account for "Able Anders"
    When I enable the account for "Able Anders"
    Then I expect the command to succeed

  Scenario: Enabling already enabled accounts return an error
    Given I create an account for the consumer "Able Anders"
    When I enable the account for "Able Anders"
    Then I expect the command to fail
    And the returned error message is "account is enabled"

  Scenario: Enabling accounts that do not exist returns an error
    Given I create an account for the consumer "Able Anders"
    And I disable the account for "Able Anders"
    When I enable the account for "Betty Burns"
    Then I expect the command to fail
    And the returned error message is "account not found"
