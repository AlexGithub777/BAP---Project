# BAP Testing
## Test 1: Login

**Component:** Authentication

**Prerequisites**: Standard login

**Steps**
* User opens app
* User inputs required information in fields (ie Email in email field, password in password field)
* User clicks login

**Test Results**
<table>
    <tr>
        <th>Input</th>
        <th>Expected Result</th>
        <th>Actual Result</th>
    </tr>
    <tr>
        <td>Username: "user1"<br>Password: "Password1!"</td>
        <td>User successfully logs in and can see a map of EIT's Campus</td>
        <td>UNTESTED</td>
    </tr>
    <tr>
        <td>Username: <BLANK><br>Password: "Password1!"</td>
        <td>Login Fails due to missing input/s</td>
        <td>UNTESTED</td>
    </tr>
    <tr>
        <td>Username: "user1"<br>Password: <BLANK></td>
        <td>Login Fails due to missing input/s</td>
        <td>UNTESTED</td>
    </tr>
    <tr>
        <td>Username: <BLANK><br>Password: <BLANK></td>
        <td>Login Fails due to missing input/s</td>
        <td>UNTESTED</td>
    </tr>
</table>

___

## Test 2: Data update

**Component:** Dashboard 

**Prerequisites**: Admin account

**Steps**
* User finds an extinguisher to change information of (map)
* User clicks update
* User inputs relevant information pertaining to the fire extinguisher
* User clicks update

**Test data**: Updated fire extinguisher information

**Expected Results**: User can update a fire extinguishers info

**Actual Results**: UNTESTED

**Test Status**: UNTESTED

**Test Status**: UNTESTED
___

## Test 3: User account to Admin account
**Steps**
* User opens user list
* User finds other user to change into admin (or searches for)
* User clicks on them, then on a drop down field saying “Set Admin”
* User can repeat steps as many times as required

**Prerequisites**: "God" account

**Test data**: None

**Expected Results**: User can update a user into an Admin

**Actual Results**: UNTESTED

**Test Status**: UNTESTED

___ 

## Test 4: View Safety Device
**Steps**
* User finds an extinguisher to open (map)
* User clicks on fire extinguisher

**Prerequisites**: User account

**Test data**: None

**Expected Results**: User can view a fire extinguishers info

**Actual Results**: UNTESTED

**Test Status**: UNTESTED

___ 

## Test 5: Sign up
**Steps**
* User opens app and clicks on "sign up"
* User fills in fields with relevant information
* User clicks sign up

**Prerequisites**: None

**Test data**: Valid email and password

**Expected Results**: User successfully creates account

**Actual Results**: UNTESTED

**Test Status**: UNTESTED
___ 

## Test 6: Deleting account
**Steps**
* User logs into account
* User clicks on account settings in app (From home: Settings -> Account Settings)
* User clicks delete account

**Prerequisites**: User account

**Test data**: None

**Expected Results**: User account is successfully deleted and logs user out

**Actual Results**: UNTESTED

**Test Status**: UNTESTED
