@query @ticket
Feature: Get Tickets

  Background: Setup resources
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

  Scenario: Can get tickets
    When I get the ticket with:
    """
    {
      "TicketID": "<TicketID>"
    }
    """
    Then I expect the request to succeed


  Scenario: Requesting tickets that do not exist returns an error
    When I get a ticket with:
    """
    {
      "TicketID": "b456"
    }
    """
    Then I expect the request to fail
    And the returned error message is:
    """
    ticket not found
    """
