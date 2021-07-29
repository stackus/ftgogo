@query @ticket
Feature: Get Tickets

  Background: Setup resources
    Given I have created a ticket for order "A123" and restaurant "Best Foods" with items
      | MenuItemID | Name       | Quantity |
      | I123       | Yummy Dish | 1        |

  Scenario: Can get tickets
    When I get the ticket for order "A123"
    Then I expect the request to succeed


  Scenario: Requesting tickets that do not exist returns an error
    When I get the ticket for order "B456"
    Then I expect the request to fail
    And the returned error message is "ticket not found"
