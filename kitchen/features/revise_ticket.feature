@command @ticket @revise
Feature: Revising Tickets

  Background: Setup Ticket
    Given I have created a ticket for order "A123" and restaurant "Best Foods" with items
      | MenuItemID | Name       | Quantity |
      | I123       | Yummy Dish | 1        |
    And I have confirmed creating a ticket for order "A123"
    And I have accepted the ticket for order "A123" will be ready in 30 minutes

  Scenario: Accepted tickets can be revised
    When I begin revising the ticket for order "A123"
    Then I expect the command to succeed

  Scenario: Tickets can be fully revised
    Given I have begun revising the ticket for order "A123"
    When I confirm revising the ticket for order "A123"
    Then I expect the command to succeed

  Scenario: Tickets being revised cannot be cancelled
    Given I have begun revising the ticket for order "A123"
    When I begin cancelling the ticket for order "A123"
    Then I expect the command to fail
    And the returned error message is "ticket state does not allow action"
