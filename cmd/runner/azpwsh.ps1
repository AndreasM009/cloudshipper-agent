Param(
    [Parameter(Mandatory=$true)]
    [string] $ScriptToRun,
    [Parameter(Mandatory=$true)]
    [string] $ArgumentsToRun,
    [Parameter(Mandatory=$true)]
    [string] $Sp,
    [Parameter(Mandatory=$true)]
    [string] $Secret,
    [Parameter(Mandatory=$true)]
    [string] $Tenant,
    [Parameter(Mandatory=$true)]
    [string] $Subscription,
    [Parameter(Mandatory=$true)]
    [string] $ArtifactsDirectory,
    [string] $WorkingDirectory = ""
)

function Get-ArgumentsTable([string] $Arguments) {
    $argArray = $Arguments.Split(" ")
    $paramNames = [System.Collections.ArrayList]@()
    $paramValues = [System.Collections.ArrayList]@()

    for ($i=0; $i -lt $argArray.Count;)
    {
        $paramNames.Add($argArray[$i]) | Out-Null

        if ($i + 1 -lt $argArray.Count)
        {
            if ($argArray[$i + 1].StartsWith("-"))
            {
                $paramValues.Add("") | Out-Null
                $i += 1
            }
            else 
            {   
                $paramValues.Add($argArray[$i + 1]) | Out-Null
                $i += 2
            }
        }
        else 
        {
            $paramValues.Add("") | Out-Null
            $i += 1
        }
    }

    $argTable = @{}

    for ($i=0; $i -lt $paramNames.Count; $i++)
    {
        $argTable[$paramNames[$i]] = $paramValues[$i]
    }

    $argTable

    Write-Host @argTable
}

Import-Module -Name Az

$argTable = Get-ArgumentsTable -Arguments $ArgumentsToRun
$pwd = ConvertTo-SecureString -String $Secret -AsPlaintext -Force
$credential = New-Object System.Management.Automation.PSCredential -ArgumentList $Sp, $pwd

Write-Host @argTable

Write-Host "Connecting to Azure"

try 
{
    Connect-AzAccount -Scope Process -Credential $credential -Tenant $Tenant -ServicePrincipal -Subscription $Subscription
} 
catch 
{
    Write-Host "Failed to connect to Azure"
    Write-Error "Connect-AzAccount failed: $_.Message"
    exit 1
}

Write-Host "Connected to Azure subscription"

$location = Join-Path -Path $ArtifactsDirectory -ChildPath $WorkingDirectory
Write-Host "Setting location to $location"
Set-Location -Path $location
# use Splatting : https://docs.microsoft.com/en-us/powershell/module/microsoft.powershell.core/about/about_splatting?view=powershell-7
&"$ScriptToRun" @argTable