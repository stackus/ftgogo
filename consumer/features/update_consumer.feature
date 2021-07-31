@command @consumer
Feature: Update Consumers

  Background: Setup a consumer
    Given I register a consumer named "Able Anders"

  Scenario: Consumers can be updated
    When I change "Able Anders" name to "Anders Able"
    Then I expect the command to succeed

  Scenario: Consumer names are changed
    Given I change "Able Anders" name to "Anders Able"
    When I request the consumer named "Anders Able"
    Then I expect the request to succeed
    And the returned consumer has the name "Anders Able"

  Scenario: Updating consumers that do not exist returns an error
    When I change "Betty Burns" name to "Burns Betty"
    Then I expect the command to fail
    And the returned error message is "consumer not found"
