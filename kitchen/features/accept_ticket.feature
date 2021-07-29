@command @ticket
Feature: Accept Tickets

  Background: Setup Ticket
    Given I have created a ticket for order "A123" at restaurant "Best Foods" with items
      | MenuItemID | Name       | Quantity |
      | I123       | Yummy Dish | 1        |
    And I have confirmed creating a ticket for order "A123"

  Scenario: Confirmed tickets can be accepted
    When I accept the ticket for order "A123" will be ready in 30 minutes
    Then I expect the command to succeed

  Scenario: Accepted tickets have the status "Accepted"
    Given I accept the ticket for order "A123" will be ready in 30 minutes
    When I get the ticket for order "A123"
    Then I expect the command to succeed
    And the returned ticket status is "Accepted"
