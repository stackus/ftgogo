@query @delivery
Feature: Get Deliveries

  Background: Setup resources
    Given a restaurant named "Best Foods" exists with address
      | Street1 | 123 Address St. |
      | City    | BigTown         |
      | State   | Colorado        |
      | Zip     | 80120           |
    And I create a delivery for order "A123" from "Best Foods" to address
      | Street1 | 123 Address St. |
      | City    | BigTown         |
      | State   | Colorado        |
      | Zip     | 80120           |

  Scenario: Can get deliveries
    When I get the delivery information for order "A123"
    Then I expect the request to succeed


  Scenario: Requesting deliveries for orders that do not exist returns an error
    When I get the delivery information for order "B456"
    Then I expect the request to fail
    And the returned error message is "delivery not found"
