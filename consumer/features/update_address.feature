@command @address
Feature: Update Consumer Address

  Background: Setup a consumer
    Given I register a consumer named "Able Anders"
    And I add an address for "Able Anders" with label "Home"
      | Street1 | 123 Address St. |
      | City    | BigTown         |
      | State   | Tristate        |
      | Zip     | 90210           |

  Scenario: Can update consumer addresses
    When I update an address for "Able Anders" with label "Home"
      | Street1 | 456 Address St. |
      | City    | BigTown         |
      | State   | Tristate        |
      | Zip     | 90210           |
    Then I expect the command to succeed

  Scenario: Updating addresses on consumers that do not exist returns an error
    When I update an address for "Betty Burns" with label "Home"
      | Street1 | 456 Address St. |
      | City    | BigTown         |
      | State   | Tristate        |
      | Zip     | 90210           |
    Then I expect the command to fail
    And the returned error message is "consumer not found"


  Scenario: Updating addresses that do not exist returns an error
    When I update an address for "Able Anders" with label "Other"
      | Street1 | 456 Address St. |
      | City    | BigTown         |
      | State   | Tristate        |
      | Zip     | 90210           |
    Then I expect the command to fail
    And the returned error message is "address with that identifier does not exist"
