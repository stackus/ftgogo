@command @delivery
Feature: Scheduling Deliveries

  Background: Setup resources
    Given a restaurant named "Best Foods" exists with address
      | Street1 | 123 Address St. |
      | City    | BigTown         |
      | State   | Tristate        |
      | Zip     | 90210           |
    And a courier exists named "Quick Courier"
    And I create a delivery for order "A123" from "Best Foods" to address
      | Street1 | 456 Address St. |
      | City    | SmallCity       |
      | State   | Tristate        |
      | Zip     | 90210           |

  Scenario: Deliveries can be scheduled with a courier
    When I schedule the delivery for order "A123"
    Then I expect the command to succeed

  Scenario: Scheduling a delivery sets the status to "SCHEDULED"
    When I schedule the delivery for order "A123"
    And I get the delivery information for order "A123"
    Then I expect the request to succeed
    And the returned delivery status is "SCHEDULED"

  Scenario: Couriers are given a two step delivery plan when scheduled with a delivery
    Given I schedule the delivery for order "A123"
    When I get the assigned courier for order "A123"
    Then I expect the request to succeed
    And the returned courier will pickup the food at address
      | Street1 | 123 Address St. |
      | City    | BigTown         |
      | State   | Tristate        |
      | Zip     | 90210           |
    And the returned courier will dropoff the food at address
      | Street1 | 456 Address St. |
      | City    | SmallCity       |
      | State   | Tristate        |
      | Zip     | 90210           |

  Scenario: Unavailable couriers are not assigned to deliveries
    Given I set the courier "Quick Courier" to be unavailable
    And I schedule the delivery for order "A123"
    When I get the assigned courier for order "A123"
    Then I expect the request to succeed
    And "Quick Courier" is not the assigned courier
