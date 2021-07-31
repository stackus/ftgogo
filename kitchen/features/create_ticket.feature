@command @ticket @create
Feature: Create Tickets

  Scenario: Can create new tickets
    When I create a ticket for order "A123" at restaurant "Best Foods" with items
      | MenuItemID | Name       | Quantity |
      | I123       | Yummy Dish | 1        |
    Then I expect the command to succeed

  Scenario: Tickets are created with a "CreatePending" status
    Given I have created a ticket for order "A123" and restaurant "Best Foods" with items
      | MenuItemID | Name       | Quantity |
      | I123       | Yummy Dish | 1        |
    When I get the ticket for order "A123"
    Then I expect the command to succeed
    And the returned ticket status is "CreatePending"

  Scenario: Tickets can be cancelled during creation
    Given I have created a ticket for order "A123" and restaurant "Best Foods" with items
      | MenuItemID | Name       | Quantity |
      | I123       | Yummy Dish | 1        |
    When I cancel creating the ticket for order "A123"
    Then I expect the command to succeed

  Scenario: Tickets cannot be cancelled after being confirmed
    Given I have created a ticket for order "A123" and restaurant "Best Foods" with items
      | MenuItemID | Name       | Quantity |
      | I123       | Yummy Dish | 1        |
    And I have confirmed creating a ticket for order "A123"
    When I cancel creating the ticket for order "A123"
    Then I expect the command to fail
    And the returned error message is "ticket state does not allow action"
