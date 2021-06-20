@command @ticket
Feature: Accept Tickets

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

  Scenario: Confirmed tickets can be accepted
    Given I have confirmed creating a ticket with:
    """
    {
      "TicketID": "<TicketID>"
    }
    """
    When I accept a ticket with:
    """
    {
      "TicketID": "<TicketID>",
      "ReadyBy": "2026-01-02T15:04:05Z"
    }
    """
    Then I expect the command to succeed

  Scenario: Confirmed tickets can be accepted
    Given I have confirmed creating a ticket with:
    """
    {
      "TicketID": "<TicketID>"
    }
    """
    When I accept the ticket with:
    """
    {
      "TicketID": "<TicketID>",
      "ReadyBy": "2026-01-02T15:04:05Z"
    }
    """
    Then I expect the command to succeed

  Scenario: Accepted tickets have the status "Accepted"
    Given I have confirmed creating a ticket with:
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
    When I get the ticket with:
    """
    {
      "TicketID": "<TicketID>"
    }
    """
    Then I expect the command to succeed
    And the returned ticket status is:
    """
    Accepted
    """
