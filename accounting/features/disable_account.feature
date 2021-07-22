@command
Feature: Disable Accounts

  Scenario: Enabled accounts can be disabled
    Given I create an account for the consumer "Able Anders"
    When I disable the account for "Able Anders"
    Then I expect the command to succeed

  Scenario: Disabling already disabled accounts return an error
    Given I create an account for the consumer "Able Anders"
    And I disable the account for "Able Anders"
    When I disable the account for "Able Anders"
    Then I expect the command to fail
    And the returned error message is "account is disabled"

  Scenario: Disabling accounts that do not exist returns an error
    Given I create an account for the consumer "Able Anders"
    When I disable the account for "Betty Burns"
    Then I expect the command to fail
    And the returned error message is "account not found"
