@command @address
Feature: Add Consumer Address

  Background: Setup a consumer
    Given I register a consumer named "Able Anders"

  Scenario: Can add delivery addresses to consumers
    When I add an address for "Able Anders" with label "Home"
      | Street1 | 123 Address St. |
      | City    | BigTown         |
      | State   | Colorado        |
      | Zip     | 80120           |
    Then I expect the command to succeed

  Scenario: Can multiple delivery addresses to consumers
    Given I add an address for "Able Anders" with label "Home"
      | Street1 | 123 Address St. |
      | City    | BigTown         |
      | State   | Colorado        |
      | Zip     | 80120           |
    When I add another address for "Able Anders" with label "Work"
      | Street1 | 123 Address St. |
      | City    | SmallCity       |
      | State   | Colorado        |
      | Zip     | 80120           |
    Then I expect the command to succeed


  Scenario: Adding an address to consumers that do not exist returns an error
    When I add an address for "Betty Burns" with label "Home"
      | Street1 | 123 Address St. |
      | City    | BigTown         |
      | State   | Colorado        |
      | Zip     | 80120           |
    Then I expect the command to fail
    And the returned error message is "consumer not found"

  Scenario: Adding an address with a duplicate id returns an error
    Given I add an address for "Able Anders" with label "Home"
      | Street1 | 123 Address St. |
      | City    | BigTown         |
      | State   | Colorado        |
      | Zip     | 80120           |
    When I add another address for "Able Anders" with label "Home"
      | Street1 | 123 Address St. |
      | City    | BigTown         |
      | State   | Colorado        |
      | Zip     | 80120           |
    Then I expect the command to fail
    And the returned error message is "address with that identifier already exists"
