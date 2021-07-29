@command @address
Feature: Remove Consumer Address

  Background: Setup a consumer
    Given I register a consumer named "Able Anders"

  Scenario: Can remove consumer addresses
    Given I add an address for "Able Anders" with label "Home"
      | Street1 | 123 Address St. |
      | City    | BigTown         |
      | State   | Tristate        |
      | Zip     | 90210           |
    When I remove an address for "Able Anders" with label "Home"
    Then I expect the command to succeed

  Scenario: Removing addresses on consumers that do not exist returns an error
    Given I add an address for "Able Anders" with label "Home"
      | Street1 | 123 Address St. |
      | City    | BigTown         |
      | State   | Tristate        |
      | Zip     | 90210           |
    When I remove an address for "Betty Burns" with label "Home"
    Then I expect the command to fail
    And the returned error message is "consumer not found"

  Scenario: Removing addresses that do not exist returns an error
    Given I add an address for "Able Anders" with label "Home"
      | Street1 | 123 Address St. |
      | City    | BigTown         |
      | State   | Tristate        |
      | Zip     | 90210           |
    When I remove an address for "Able Anders" with label "Other"
    Then I expect the command to fail
    And the returned error message is "address with that identifier does not exist"
