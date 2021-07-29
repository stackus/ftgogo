@command @order @revise
Feature: Order Revision

  Background: Setup Resources
    Given I have initialized the restaurant "Best Foods"
    And I have submitted an order to "Best Foods" from "Able Anders"

  Scenario: Approved orders may be revised
    Given I have approved the order to "Best Foods" from "Able Anders" with ticket "T123"
    When I begin to revise the order to "Best Foods" from "Able Anders"
    Then I expect the command to succeed
    And expect the order to "Best Foods" from "Able Anders" is "RevisionPending"

  Scenario: Rejected orders cannot be revised
    Given I have rejected the order to "Best Foods" from "Able Anders"
    When I begin to revise the order to "Best Foods" from "Able Anders"
    Then I expect the command to fail
    And the returned error message is "order state does not allow action"

  Scenario: Can confirm revision for orders pending revision
    Given I have approved the order to "Best Foods" from "Able Anders" with ticket "T123"
    And I have begun to revise the order to "Best Foods" from "Able Anders"
    When I confirm revising the order to "Best Foods" from "Able Anders"
    Then I expect the command to succeed
    And expect the order to "Best Foods" from "Able Anders" is "Approved"

  Scenario: Can undo revision for orders pending revision
    Given I have approved the order to "Best Foods" from "Able Anders" with ticket "T123"
    And I have begun to revise the order to "Best Foods" from "Able Anders"
    When I undo revising the order to "Best Foods" from "Able Anders"
    Then I expect the command to succeed
    And expect the order to "Best Foods" from "Able Anders" is "Approved"
