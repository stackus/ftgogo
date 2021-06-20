@command @ticket @cancel
Feature: Cancelling Tickets

  Background: Setup Ticket
    Given I create a ticket with:
    """
    {
      "OrderID": "a123",
      "RestaurantID": "a123",
      "LineItems": [
        {
          "MenuItemID": "a123",
          "Name": "TestMenuItem",
          "Quantity": 1
        }
      ]
    }
    """
    And I have confirmed creating a ticket with:
    """
    {
      "TicketID": "<TicketID>"
    }
    """
    And I have accepted the ticket with:
    """
    {
      "TicketID": "<TicketID>",
      "ReadyBy": "2026-01-02T15:04:05Z"
    }
    """

  Scenario: Accepted tickets can be cancelled
    When I begin cancelling the ticket with:
    """
    {
      "TicketID": "<TicketID>"
    }
    """
    Then I expect the command to succeed

  Scenario: Tickets can be fully cancelled
    Given I have begun cancelling the ticket with:
    """
    {
      "TicketID": "<TicketID>"
    }
    """
    When I confirm cancelling the ticket with:
    """
    {
      "TicketID": "<TicketID>"
    }
    """
    Then I expect the command to succeed

  Scenario: Tickets being cancelled cannot be revised
    Given I have begun cancelling the ticket with:
    """
    {
      "TicketID": "<TicketID>"
    }
    """
    When I begin revising the ticket with:
    """
    {
      "TicketID": "<TicketID>"
    }
    """
    Then I expect the command to fail
    And the returned error message is:
    """
    ticket state does not allow action
    """
