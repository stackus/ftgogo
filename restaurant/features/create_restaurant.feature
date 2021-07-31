Feature: Restaurant Creation

  Scenario: Can create new restaurants
    When I create the restaurant "Best Foods"
    Then I expect the command to succeed
