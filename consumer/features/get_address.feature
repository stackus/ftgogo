@query @address
Feature: Get Consumer Address

  Background: Setup a consumer
    Given I register a consumer named "Able Anders"

  Scenario: Can get a consumers address
    Given I add an address for "Able Anders" with label "Home"
      | Street1 | 123 Address St. |
      | City    | BigTown         |
      | State   | Colorado        |
      | Zip     | 80120           |
    When I request the "Home" address for "Able Anders"
    Then I expect the request to succeed
    And the returned address to match
      | Street1 | 123 Address St. |
      | City    | BigTown         |
      | State   | Colorado        |
      | Zip     | 80120           |

  Scenario: Getting an address that doesn't exist returns an error
    Given I add an address for "Able Anders" with label "Home"
      | Street1 | 123 Address St. |
      | City    | BigTown         |
      | State   | Colorado        |
      | Zip     | 80120           |
    When I request the "Other" address for "Able Anders"
    Then I expect the request to fail
    And the returned error message is "an address with that identifier does not exist"

  Scenario: Getting an address for a consumer that doesn't exist returns an error
    Given I add an address for "Able Anders" with label "Home"
      | Street1 | 123 Address St. |
      | City    | BigTown         |
      | State   | Colorado        |
      | Zip     | 80120           |
    When I request the "Home" address for "Betty Burns"
    Then I expect the request to fail
    And the returned error message is "consumer not found"
