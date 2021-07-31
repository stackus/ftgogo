@query @delivery
Feature: Get Deliveries

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

  Scenario: Can get deliveries
    When I get the delivery information for order "A123"
    Then I expect the request to succeed


  Scenario: Requesting deliveries for orders that do not exist returns an error
    When I get the delivery information for order "B456"
    Then I expect the request to fail
    And the returned error message is "delivery not found"
