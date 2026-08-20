% Stone-skipping simulation
% The calculations follow the original project script; comments and spacing
% have been clarified for display on the project page.

z = 50;
image_size_x = 6 * z;
image_size_y = 2;

theta = pi / 200;      % Stone angle above the water surface
size_of_stone = 5;     % Stone length
mass = 1;              % Stone mass
G = 10;                % Gravitational acceleration

% Initial position and velocity
x = -image_size_x;
y = 0.5;
vx = 3 * z;
vy = 0.5;
dt = 0.005;            % Euler time step

rho = 1;               % Water density
C1 = 0.1;              % Lift coefficient
C2 = 0.3;              % Drag coefficient

vx_data = vx;
vy_data = vy;
t_data = 0;
data_index = 1;

% Flat water surface
[X, Y] = meshgrid(...
    linspace(-image_size_x, image_size_x, 25), ...
    linspace(-image_size_x, image_size_x, 25));
Z = zeros(size(X));

% Expanding circles mark each water contact
circle_theta = linspace(0, 2 * pi, 100);
circle = 0;
radius = 0;
circle_x = 0;
number_of_circle = 0;

while y > -image_size_y
    if y > 0
        % In the air, only gravity changes the vertical velocity.
        vy = vy - G * dt;

        if circle == 1
            circle = 0;
        end
    else
        % Record the first time the stone enters the water on each skip.
        if circle == 0
            number_of_circle = number_of_circle + 1;
            circle_x(number_of_circle) = x;
            radius(number_of_circle) = 0;
            circle = 1;
        end

        if y + size_of_stone * sin(theta) > 0
            % The stone is only partly submerged. Lift and drag depend on
            % speed, submerged length, and the stone's angle.
            Sim = min(abs(y), size_of_stone) / sin(theta);
            vvy = -G + 0.5 * rho / mass * (vy^2 + vx^2) * Sim * ...
                (C1 * cos(theta) - C2 * sin(theta));
            vvx = -0.5 * rho / mass * (vy^2 + vx^2) * Sim * ...
                (C1 * sin(theta) + C2 * cos(theta));
        else
            % The stone is fully submerged.
            vvx = -C2 * vx^2;
            vvy = -G + C2 * vy^2;
        end

        vx = vx + vvx * dt;
        vy = vy + vvy * dt;
    end

    x = x + vx * dt;
    y = y + vy * dt;

    % Draw the stone.
    plot3(...
        [x, x + size_of_stone * cos(theta)], ...
        [0, 0], ...
        [y, y + size_of_stone * sin(theta)], ...
        'LineWidth', 4, ...
        'Color', 'r');
    hold on;
    axis([...
        -image_size_x, image_size_x, ...
        -image_size_x, image_size_x, ...
        -image_size_y, image_size_y]);

    % Draw a growing ripple at every impact point.
    if number_of_circle > 0
        for i = 1:number_of_circle
            plot3(...
                circle_x(i) + radius(i) * cos(circle_theta), ...
                radius(i) * sin(circle_theta), ...
                zeros(size(circle_theta)));
            radius(i) = radius(i) + 1.5;
        end
    end

    mesh(X, Y, Z);
    Z = zeros(size(X));

    xlabel('distance');
    zlabel('height');
    title('Trajectory of stone');
    hold off;

    % Save velocity history for optional diagnostic plots.
    t_data = [t_data, data_index * dt]; %#ok<AGROW>
    vx_data = [vx_data, vx]; %#ok<AGROW>
    vy_data = [vy_data, vy]; %#ok<AGROW>
    data_index = data_index + 1;

    drawnow;
end
