@command @ticket @cancel
Feature: Cancelling Tickets

  Background: Setup Ticket
    Given I have created a ticket for order "A123" at restaurant "Best Foods" with items
      | MenuItemID | Name       | Quantity |
      | I123       | Yummy Dish | 1        |
    And I have confirmed creating a ticket for order "A123"
    And I have accepted the ticket for order "A123" will be ready in 30 minutes

  Scenario: Accepted tickets can be cancelled
    When I begin cancelling the ticket for order "A123"
    Then I expect the command to succeed

  Scenario: Tickets can be fully cancelled
    Given I have begun cancelling the ticket for order "A123"
    When I confirm cancelling the ticket for order "A123"
    Then I expect the command to succeed

  Scenario: Tickets being cancelled cannot be revised
    Given I have begun cancelling the ticket for order "A123"
    When I begin revising the ticket for order "A123"
    Then I expect the command to fail
    And the returned error message is "ticket state does not allow action"
