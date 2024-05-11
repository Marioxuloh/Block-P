ssh -i "garden-keys.pem" admin@15.188.167.146

scp -i "garden-keys.pem" mysql-community-server_8.4.0-1debian12_amd64.deb admin@ec2-51-44-11-123.eu-west-3.compute.amazonaws.com:/home/admin/dependencies
mysql -h garden-database.cn24a22i2d47.eu-west-3.rds.amazonaws.com -P 3306 -u gardenAdmin -p BOTB2CZnKX8hK4vzlSEB

scp -i "garden-keys.pem" /Users/marioxuloh/Documents/VSCode/garden-backend/bin/Release/net8.0/garden-backend.dll admin@15.188.167.146:/home/admin/Garden-Backend
scp -i "garden-keys.pem" -r /Users/marioxuloh/Documents/VSCode/garden-backend/bin/Release/net8.0/publish/ admin@15.188.167.146:/home/admin/Garden-Backend


aqui los servicios que ejecutaran todo en el ec2 /etc/systemd/system/

aws acm get-certificate --certificate-arn arn:aws:acm:us-east-1:211125788925:certificate/8d75ed81-48fd-4f07-aa37-c69355fe7254 --query 'Certificate' --output text > cert.pem

      /*"HttpsInLineCertFile": {
        "Url": "https://0.0.0.0:5041",
        "Certificate": {
          "Path": "/home/admin/Garden-Backend/certificate.pfx",
          "Password": "8764"
        }
      }*/

      
builder.Services.AddHttpsRedirection(options =>
{
    options.RedirectStatusCode = StatusCodes.Status307TemporaryRedirect;
    options.HttpsPort = 5041;
});

var app = builder.Build();

app.UseStaticFiles();
app.UseAuthentication();
app.UseAuthorization();
app.UseCors("AllowAnyOrigin");
app.UseHttpsRedirection();
app.MapControllers();
app.Run();