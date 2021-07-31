@query @consumer
Feature: Get Consumer

  Background: Setup a consumer
    Given I register a consumer named "Able Anders"

  Scenario: Can get consumers
    When I request the consumer named "Able Anders"
    Then I expect the request to succeed

  Scenario: Asking for a consumer that does not exist returns an error
    When I request the consumer named "Betty Burns"
    Then I expect the request to fail
    And the returned error message is "consumer not found"

  Scenario: Consumer is returned with labeled addresses
    Given I add an address for "Able Anders" with label "Home"
      | Street1 | 123 Address St. |
      | City    | BigTown         |
      | State   | Tristate        |
      | Zip     | 90210           |
    When I request the consumer named "Able Anders"
    Then I expect the request to succeed
    And the returned consumer has an address with label "Home"

  Scenario: Consumer is returned with all addresses
    Given I add an address for "Able Anders" with label "Home"
      | Street1 | 123 Address St. |
      | City    | BigTown         |
      | State   | Tristate        |
      | Zip     | 90210           |
    Given I add an address for "Able Anders" with label "Work"
      | Street1 | 123 Address St. |
      | City    | SmallCity       |
      | State   | Tristate        |
      | Zip     | 90210           |
    When I request the consumer named "Able Anders"
    Then I expect the request to succeed
    And the returned consumer has 2 addresses
