using Consul;

namespace GrpcForConsul
{
    public static class RegisterConsul
    {
        public static IServiceCollection ConsulRegister<TServer>(this IServiceCollection services)
        {
            var client = services.BuildServiceProvider().GetService<IConsulClient>();
            var registerID = $"core.grpc.server.{nameof(TServer)}";
            client.Agent.ServiceDeregister(registerID).Wait();
            client.Agent.ServiceRegister(new AgentServiceRegistration()
            {
                ID = registerID,
                Name = typeof(TServer).Name,
                Address = "localhost",
                Port = 5030,
                Check = new AgentServiceCheck
                {
                    TCP = "192.168.3.3:5030",
                    Status = HealthStatus.Passing,
                    DeregisterCriticalServiceAfter = TimeSpan.FromSeconds(10),
                    Interval = TimeSpan.FromSeconds(10),
                    Timeout = TimeSpan.FromSeconds(5)
                },
                Tags = new string[] { "gRpc" }
            }).Wait();
            return services;
        }
    }
}
