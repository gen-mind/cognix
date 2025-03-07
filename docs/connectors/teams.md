#### Microsoft Teams Teams 
This is the description for the Microsoft Teams Team and private and group chats connector

## Setup 
From CogniX UI navigate to connectors and create a new connector
Choose `Teams`
At step 2:
- Choose a name, it's just a description
- Fill the "Connector Specific Configration" with the json below filled with the corect data
- Refresh frequency in seconds is the delta of time that CogniX will use to start a new scan on your connected data source
- Connector credential, fill with a random number, it's not used

```json
{
  "team": "",
  "analyze_chats": true, 
  "channels": [],
  "files": {
    "folder": "",
    "recursive": true
  },
  "token": {
    "access_token": "",
    "expiry": "",
    "refresh_token": "",
    "token_type": ""
  }
}
```

- team : name of team for analyzing. teams channel will not be analyzed if empty string. 
- channels : name of channels for analyzing ( all channels will be analyzed if array is empty )
- analyze_chats : true for analyze all private and group chats.
- token : OAuth token for access to ```ms teams```
- files : 
  - folder : optional, folder name for scanning
  - recursive :  false - scan only given folder , true - scan nested folders


Since the UI is still under construction you'll need to do some manual steps to get the OneDrive token.
This process will be automated with the UI evolution

Paste in your browser the following link if you are running CogniX on your private Docker deployment
```js
    http://localhost:8080/api/oauth/microsoft/auth_url?redirect_url=http://localhost:8080
```

If you are using CogniX from [rag.cognix.ch](https://rag.cognix.ch)
```js
    https://rag.cognix.ch/api/oauth/microsoft/auth_url?redirect_url=https://rag.cognix.ch
```

once you paste the link above in the browser you will get a string. copy link <br/>
you will get something similar to the json below:<br/>

```json
https://login.microsoftonline.com/common/oauth2/v2.0/authorize?client_id=<id>>&scope=offline_access Files.Read.All Sites.ReadWrite.All&response_type=code&redirect_uri=http://localhost:8080/api/oauth/microsoft/callback
```

paste the link as described above in a new browser window <br/>. 
Sign in using you microsoft account and grant permission to CogniX<br/>
There's a checkbox you need to mark "Consent on behalf of your company"<br/>
Click Accept <br/>

You will be prompted with another json similar to the one described above<br/>
This is a overly long string, though.
Be prepared with an editor so that you will be able to handle correctly the copy/paste operations
Copy the token from the response<br/>
The property named "access_token", "expiry": refresh_token": "", "token_type" and paste in the json provided above <br/>
It might be a bit complex because access token and refresh token are very long string
Make sure to copy them properly
The token that you will receive will look like the sample below

```json
{
  "status": 200,
  "data": {
    "id": "",
    "email": "",
    "name": "",
    "given_name": "",
    "family_name": "",
    "access_token": "",
    "refresh_token": "",
    "token": {
      "access_token": "",
      "token_type": "",
      "refresh_token": "",
      "expiry": ""
    }
  }
}
```

Now you have a json filled with all the values CogniX needs.<br/>
Paste it into the connector specific configuration <br/>

**Refresh frequency** is in second it tells CogniX every each seconds it need to rescan the source.
Make it at least 86400 (one day in seconds) <br/>
**connector credentials**
not used add a number

## Microsoft Configuration
If you don't have a M365 account, you need to set it up. Instructions [here](https://learn.microsoft.com/en-us/microsoft-365/admin/simplified-signup/signup-business-standard?view=o365-worldwide#sign-up-for-microsoft-365-for-business)

Once you have a M365 account
Go to our Azure subscription and create a new app registration

