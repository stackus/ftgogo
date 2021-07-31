@command
Feature: Create Account

  Scenario: Create a new account
    When I create an account for the consumer "Able Anders"
    Then I expect the command to succeed

  Scenario: Creating a duplicate account returns an error
    Given I create an account for the consumer "Able Anders"
    When I create an account for the consumer "Able Anders"
    Then I expect the command to fail
    And the returned error message is "account already exists"
