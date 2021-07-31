@command @delivery
Feature: Create Deliveries

  Background: Setup a restaurant
    Given a restaurant named "Best Foods" exists with address
      | Street1 | 123 Address St. |
      | City    | BigTown         |
      | State   | Tristate        |
      | Zip     | 90210           |

  Scenario: Can create deliveries
    When I create a delivery for order "A123" from "Best Foods" to address
      | Street1 | 123 Address St. |
      | City    | BigTown         |
      | State   | Tristate        |
      | Zip     | 90210           |
    Then I expect the command to succeed

  Scenario: Deliveries are created with a "PENDING" status
    When I create a delivery for order "A123" from "Best Foods" to address
      | Street1 | 123 Address St. |
      | City    | BigTown         |
      | State   | Tristate        |
      | Zip     | 90210           |
    And I request the delivery information for order "A123"
    Then I expect the command to succeed
    And the returned delivery status is "PENDING"

  Scenario: Creating deliveries for restaurants that do not exist returns an error
    When I create a delivery for order "A123" from "Other Foods" to address
      | Street1 | 123 Address St. |
      | City    | BigTown         |
      | State   | Tristate        |
      | Zip     | 90210           |
    Then I expect the command to fail
    And the returned error message is "restaurant not found"

  Scenario: Creating duplicate deliveries for an order returns an error
    Given I create a delivery for order "A123" from "Best Foods" to address
      | Street1 | 123 Address St. |
      | City    | BigTown         |
      | State   | Tristate        |
      | Zip     | 90210           |
    When I create another delivery for order "A123" from "Best Foods" to address
      | Street1 | 123 Address St. |
      | City    | BigTown         |
      | State   | Tristate        |
      | Zip     | 90210           |
    Then I expect the command to fail
    And the returned error message is "delivery already exists"
