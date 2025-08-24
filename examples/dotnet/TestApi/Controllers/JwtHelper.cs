using System;
using System.Collections.Generic;
using System.IdentityModel.Tokens.Jwt;
using System.Security.Claims;
using System.Text;
using Microsoft.IdentityModel.Tokens;
using Newtonsoft.Json;
using System.ComponentModel.DataAnnotations;
public static class JwtHelper
{
    public static string GenerateJwt<T>(T data, string secret, TimeSpan expire)
    {
        // 1. Convert struct/class thành dictionary (claims)
        var json = JsonConvert.SerializeObject(data);
        var dict = JsonConvert.DeserializeObject<Dictionary<string, object>>(json);

        var claims = new List<Claim>();
        if (dict != null)
        {
            foreach (var kv in dict)
            {
                if (kv.Value != null)
                {
                    claims.Add(new Claim(kv.Key, kv.Value.ToString() ?? ""));
                }
            }
        }

        // 2. Thêm Expire claim
        var key = new SymmetricSecurityKey(Encoding.UTF8.GetBytes(secret));
        var creds = new SigningCredentials(key, SecurityAlgorithms.HmacSha256);

        var token = new JwtSecurityToken(
            claims: claims,
            expires: DateTime.UtcNow.Add(expire),
            signingCredentials: creds
        );

        // 3. Trả về JWT string
        return new JwtSecurityTokenHandler().WriteToken(token);
    }
}
