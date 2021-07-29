@command @delivery
Feature: Cancel Deliveries

  Background: Setup resources
    Given a restaurant named "Best Foods" exists with address
      | Street1 | 123 Address St. |
      | City    | BigTown         |
      | State   | Tristate        |
      | Zip     | 90210           |
    And I create a delivery for order "A123" from "Best Foods" to address
      | Street1 | 123 Address St. |
      | City    | BigTown         |
      | State   | Tristate        |
      | Zip     | 90210           |


  Scenario: Cancel existing deliveries
    When I cancel delivery for order "A123"
    Then I expect the command to succeed

  Scenario: Canceling deliveries that do not exist returns an error
    When I cancel delivery for order "B456"
    Then I expect the command to fail
    And the returned error message is "delivery not found"
