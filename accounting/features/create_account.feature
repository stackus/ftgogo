Feature: Create Account

  Scenario: Create a new account
    When I create an account with:
      | ConsumerID | a123        |
      | Name       | TestAccount |
    Then I expect the command to succeed

  Scenario: Creating a duplicate account returns an error
    Given I create an account with:
      | ConsumerID | a123        |
      | Name       | TestAccount |
    When I create an account with:
      | ConsumerID | a123       |
      | Name       | NewAccount |
    Then I expect the command to fail
    And the returned error message is:
    """
    account already exists
    """
