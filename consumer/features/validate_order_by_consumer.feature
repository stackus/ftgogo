@command @consumer @order
Feature: Validate Orders By Consumer

  Background: Setup a consumer
    Given I register a consumer named "Able Anders"

  Scenario: Can validate orders for consumers
    When I validate an order for "Able Anders"
      | OrderID | A123  |
      | Total   | $9.99 |
    Then I expect the command to succeed

  Scenario: Cannot validate orders for consumers that do not exist
    When I validate an order for "Betty Burns"
      | OrderID | A123  |
      | Total   | $9.99 |
    Then I expect the command to fail
    And the returned error message is "consumer not found"
