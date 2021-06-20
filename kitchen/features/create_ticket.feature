@command @ticket @create
Feature: Create Tickets

  Scenario: Can create new tickets
    When I create a ticket with:
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
    Then I expect the command to succeed

  Scenario: Tickets are created with a "CreatePending" status
    Given I have created a ticket with:
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
    When I get the ticket with:
    """
    {
      "TicketID": "<TicketID>"
    }
    """
    Then I expect the command to succeed
    And the returned ticket status is:
    """
    CreatePending
    """

  Scenario: Tickets can be cancelled during creation
    Given I have created a ticket with:
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
    When I cancel creating the ticket with:
    """
    {
      "TicketID": "<TicketID>"
    }
    """
    Then I expect the command to succeed

  Scenario: Tickets cannot be cancelled after being confirmed
    Given I have created a ticket with:
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
    And I have confirmed creating the ticket with:
    """
    {
      "TicketID": "<TicketID>"
    }
    """
    When I cancel creating the ticket with:
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
